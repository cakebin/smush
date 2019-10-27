import { SortDirection, SortEvent } from './../../../modules/common-ux/common-ux.view-models';
import { faTrash, faPencilAlt } from '@fortawesome/free-solid-svg-icons';
import { Component, OnInit, ViewChildren, QueryList } from '@angular/core';

import { ITagViewModel, IServerResponse } from 'client/app/app.view-models';
import { CommonUxService } from 'client/app/modules/common-ux/common-ux.service';
import { HeaderViewModel, ISortEvent } from 'client/app/modules/common-ux/common-ux.view-models';
import { TagManagementService } from 'client/app/modules/tag-management/tag-management.service';
import { SortableTableHeaderComponent } from 'client/app/modules/common-ux/components/sortable-table-header/sortable-table-header.component';


@Component({
  selector: 'admin-tag',
  templateUrl: './admin-tag.component.html',
  styleUrls: ['./admin-tag.component.css']
})
export class AdminTagComponent implements OnInit {
  public headerLabels: HeaderViewModel[] = [
    new HeaderViewModel('tagId', 'ID', '100px'),
    new HeaderViewModel('tagName', 'Name', '200px'),
  ];
  @ViewChildren(SortableTableHeaderComponent) headerComponents: QueryList<SortableTableHeaderComponent>;

  public tags: ITagViewModel[] = [];
  public newTag: ITagViewModel = {} as ITagViewModel;

  public sortedTags: ITagViewModel[];
  public sortColumnName: string = '';
  public sortColumnDirection: SortDirection = '';
  public isInitialLoad: boolean = true;

  public faTrash = faTrash;
  public faPencilAlt = faPencilAlt;

  constructor(
    private tagService: TagManagementService,
    private commonUxService: CommonUxService,
  ) { }

  ngOnInit() {
    this.tagService.cachedTags.subscribe(
      (res: ITagViewModel[]) => {
      if (res && res.length) {
        this.sortedTags = res;
        this.tags = res;

        if (this.isInitialLoad) {
          this._initialSort();
          this.isInitialLoad = false;
        } else {
          this._sort();
        }
      }
    });
  }

  public createTag(): void {
    if (!this.newTag.tagName) {
      // This shouldn't happen unless someone manually re-enables the create button
      this.commonUxService.showWarningToast('Please specify text for your new tag.');
      return;
    }
    this.tagService.createTag(this.newTag).subscribe(
      (res: IServerResponse) => {
        if (res.success) {
          this.newTag = {} as ITagViewModel;
          this.commonUxService.showStandardToast('Tag created!');
        } else {
          this.commonUxService.showDangerToast('Unable to create tag.');
        }
      },
      error => {
        this.commonUxService.showDangerToast('Unable to create tag.');
        console.error(error);
      }
    );
  }

  public onSort({column, direction}: ISortEvent) {
    // Resetting all headers. This needs to be done in a parent, no way around it
    if (this.headerComponents) {
      this.headerComponents.forEach(header => {
        if (header.propertyName !== column) {
          header.clearDirection();
        }
      });
    }

    // Sorting items
    if (direction === '') {
      this.sortColumnName = '';
      this.sortColumnDirection = '';
      this.sortedTags = this.tags;
    } else {
      this.sortColumnName = column;
      this.sortColumnDirection = direction;
      this.sortedTags = [...this.tags].sort((a, b) => {
        const res = this.commonUxService.compare(a[column], b[column]);
        return direction === 'asc' ? res: -res;
      });
    }
  }

  /*------------------------
       Private helpers
  -------------------------*/
  private _initialSort(): void {
    this.onSort(new SortEvent('tagId', 'asc'));
  }

  private _sort(): void {
    this.onSort(new SortEvent(this.sortColumnName, this.sortColumnDirection));
  }
}
