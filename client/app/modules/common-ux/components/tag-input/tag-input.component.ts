import { Component, Input, Output, EventEmitter, OnInit} from '@angular/core';
import { ITagViewModel } from 'client/app/app.view-models';


@Component({
  selector: 'common-ux-tag-input',
  templateUrl: 'tag-input.component.html',
  styleUrls: ['../../common-ux.css'],
})
export class TagInputComponent implements OnInit {
    @Input() tags: ITagViewModel[] = [];
    @Input() size: 'sm' | '';
    @Input() placement: 'right' | 'top' | 'left' | 'bottom' = 'top';
    @Input() popoverClass: string = '';
    @Input() selectedTags: ITagViewModel[] = [];
    @Output() selectedTagsChange = new EventEmitter<ITagViewModel[]>();

    public newTag: ITagViewModel;

    ngOnInit() {
      if (this.popoverClass) {
        this.popoverClass = 'tag-popover ' + this.popoverClass;
      } else {
        this.popoverClass = 'tag-popover';
      }
    }

    public removeTag(tag: ITagViewModel): void {
      const tagIndex: number = this.selectedTags.findIndex(t => t.tagId === tag.tagId);
      this.selectedTags.splice(tagIndex, 1);
    }
    public onSelectTag(event: ITagViewModel): void {
      if (event != null) {
        if (!this.selectedTags.find(t => t.tagId === event.tagId)) {
          this.selectedTags.push(event);
        }
      }
    }
}
