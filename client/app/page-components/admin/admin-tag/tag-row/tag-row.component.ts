import { Component, OnInit, Input } from '@angular/core';
import { faPencilAlt, faTrash } from '@fortawesome/free-solid-svg-icons';
import { NgbPopover } from '@ng-bootstrap/ng-bootstrap';

import { CommonUxService } from 'client/app/modules/common-ux/common-ux.service';
import { TagManagementService } from 'client/app/modules/tag-management/tag-management.service';
import { ITagViewModel, IServerResponse } from 'client/app/app.view-models';


@Component({
  selector: '[tag-row]',
  templateUrl: './tag-row.component.html',
  styleUrls: ['./tag-row.component.css']
})
export class TagRowComponent implements OnInit {
  @Input() tag: ITagViewModel = {} as ITagViewModel;

  public editedTag: ITagViewModel = {} as ITagViewModel;
  public warnings: string[] = [];

  public faTrash = faTrash;
  public faPencilAlt = faPencilAlt;


  constructor(
    private tagService: TagManagementService,
    private commonUxService: CommonUxService,
  ) { }

  ngOnInit() {
    this.tag.editMode = false;
  }

  public editTag(originalTag: ITagViewModel): void {
    originalTag.editMode = true;

    this.editedTag = {
      tagId: originalTag.tagId,
      tagName: originalTag.tagName,
    } as ITagViewModel;
  }

  public deleteTag(tag: ITagViewModel): void {
    this.commonUxService.openConfirmModal(`Removing tag ${tag.tagName}`, 'Delete tag', false, 'Nuke it')
      .then(confirm => { this.tagService.deleteTag(tag); }, reject => { /* Do nothing */ });
  }

  public saveChanges(): void {
    this.tagService.updateTag(this.editedTag).subscribe(
      (res: IServerResponse) => {
        if (res) {
          console.log('returned tag', res)
          this.tag = res.data.tag;
          this.editedTag = res.data.tag;
          this.resetState();
        }
      }
    )
  }

  public resetState(): void {
    this.editedTag = {} as ITagViewModel;
    this.tag.editMode = false;
  }

  public validateTag(): boolean {
    this.warnings = [];

    if (!this.editedTag.tagName) {
      this.warnings.push('Tag Name required.');
    }

    if (this.warnings.length) {
      return false;
    } else {
      return true;
    }
  }

  public openWarningPopover(popover: NgbPopover) {
    if (!this.validateTag()) {
      popover.popoverTitle = 'Invalid tag';
      popover.ngbPopover = this.warnings.join(' ');
      popover.open();
    }
  }
}
