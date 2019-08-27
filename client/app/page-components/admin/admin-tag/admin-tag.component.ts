import { Component, OnInit } from '@angular/core';

import { ITagViewModel, IServerResponse } from 'client/app/app.view-models';
import { CommonUxService } from 'client/app/modules/common-ux/common-ux.service';
import { TagManagementService } from 'client/app/modules/tag-management/tag-management.service';


@Component({
  selector: 'admin-tag',
  templateUrl: './admin-tag.component.html',
  styleUrls: ['./admin-tag.component.css']
})
export class AdminTagComponent implements OnInit {
  public tags: ITagViewModel[] = [];
  public newTag: ITagViewModel = {} as ITagViewModel;

  constructor(
    private tagService: TagManagementService,
    private commonUxService: CommonUxService,
  ) { }

  ngOnInit() {
    this.tagService.cachedTags.subscribe(
      (res: ITagViewModel[]) => {
      if (res && res.length) {
        this.tags = res;
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
          this.commonUxService.showSuccessToast('Tag created!');
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
}
