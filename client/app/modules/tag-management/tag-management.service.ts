import { Injectable, Inject  } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { publish, refCount, tap } from 'rxjs/operators';
import { ICharacterViewModel, IServerResponse, ITagViewModel, TagViewModel } from '../../app.view-models';

@Injectable()
export class TagManagementService {

    public cachedTags: BehaviorSubject<ITagViewModel[]> = new BehaviorSubject<ITagViewModel[]>(null);

    constructor(
        private httpClient: HttpClient,
        @Inject('TagApiUrl') private apiUrl: string,
    ) {
    }
    public loadAllTags(): void {
        const fakeTags: ITagViewModel[] = [
            new TagViewModel(1, 'Laggy match'),
            new TagViewModel(2, 'Teabagging opponent'),
            new TagViewModel(3, 'Camping opponent'),
            new TagViewModel(4, 'Spamming opponent'),
            new TagViewModel(5, 'Got homie stock'),
            new TagViewModel(6, 'Gave homie stock'),
            new TagViewModel(7, 'Rematched'),
        ];
        this.cachedTags.next(fakeTags);
        this.cachedTags.pipe(
            publish(),
            refCount()
        );

        /*
        this.httpClient.get<IServerResponse>(`${this.apiUrl}/getall`).subscribe(
            (res: IServerResponse) => {
                if (res && res.data && res.data.tags) {
                    this.cachedTags.next(res.data.tags);
                    this.cachedTags.pipe(
                        publish(),
                        refCount()
                    );
                }
            }
        );
        */
    }
    public createTag(tag: ITagViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, tag).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.tag) {
                    const allTags: ITagViewModel[] = this.cachedTags.value;
                    allTags.push(res.data.tag);
                    this.cachedTags.next(allTags);
                    this.cachedTags.pipe(
                        publish(),
                        refCount()
                    );
                }
            })
        );
    }
    public updateTag(updatedTag: ITagViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/update`, updatedTag).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.tag) {
                    // Replace old tag with updated tag in a copy of cached tags
                    const updatedTagFromServer = res.data.tag;
                    const allTags: ITagViewModel[] = this.cachedTags.value;
                    const tagIndex: number = allTags.findIndex(
                        c => c.tagId === updatedTagFromServer.tagId);
                    Object.assign(allTags[tagIndex], updatedTagFromServer);

                    // Overwrite cache with updated copy
                    this.cachedTags.next(allTags);
                    this.cachedTags.pipe(
                        publish(),
                        refCount()
                    );
                }
            })
        );
    }
}
