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
        this.httpClient.get<IServerResponse>(`${this.apiUrl}/getall`).subscribe(
            (res: IServerResponse) => {
                if (res && res.data && res.data.tags) {
                    this._updateCachedTags(res.data.tags);
                }
            }
        );
    }

    public createTag(tag: ITagViewModel): Observable<{}> {
        return this.httpClient.post(`${this.apiUrl}/create`, tag).pipe(
            tap((res: IServerResponse) => {
                if (res && res.data && res.data.tag) {
                    const allTags: ITagViewModel[] = this.cachedTags.value;
                    allTags.push(res.data.tag);
                    this._updateCachedTags(allTags);
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
                    this._updateCachedTags(allTags);
                }
            })
        );
    }

   public deleteTag(tag: ITagViewModel): void {
       this.httpClient.post(`${this.apiUrl}/delete`, tag).pipe(
           tap((res: IServerResponse) => {
             if (res && res.success) {
               const allTags: ITagViewModel[] = this.cachedTags.value;
               const index = allTags.findIndex(t => t.tagId === tag.tagId);
               allTags.splice(index, 1);
               this._updateCachedTags(allTags);
             }
           })
       ).subscribe();
   }

    private _updateCachedTags(tags: ITagViewModel[]): void {
        tags.sort((a, b) => {
            if (a.tagName < b.tagName) {
                return -1;
            }
            if (a.tagName > b.tagName) {
                return 1;
            }
            if (a.tagName === b.tagName) {
                return 0;
            }
        });

        this.cachedTags.next(tags);
        this.cachedTags.pipe(
            publish(),
            refCount()
        );
    }
}
