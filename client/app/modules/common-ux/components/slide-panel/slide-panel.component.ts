import { Component, Input, OnInit, HostBinding } from '@angular/core';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { SlidePanelService } from './slide-panel.service';

type SlideDirection = 'left' | 'right';

// A fleshed out and re-designed version of the following experiment:
// https://medium.com/@asm/animated-slide-panel-with-angular-e985ad646f9

@Component({
  selector: 'common-ux-slide-panel',
  styleUrls: [ './slide-panel.component.css' ],
  templateUrl: './slide-panel.component.html',
  animations: [
    trigger('slide', [
      state('visibleLeft', style({ transform: 'translateX(0)' })),
      state('hiddenLeft', style({ transform: 'translateX(50%)' })),
      state('visibleRight', style({ transform: 'translateX(-50%)' })),
      state('hiddenRight', style({ transform: 'translateX(0)' })),
      transition('* => *', animate(300))
    ])
  ]
})
export class SlidePanelComponent implements OnInit {
  @Input() slideDirection: SlideDirection = 'left';

  public set paneVisible(isVisible: boolean) {
    this._paneVisible = isVisible;
    if (isVisible) {
      if (this.slideDirection === 'left') {
        this.paneAnimationState = 'visibleLeft';
      } else if (this.slideDirection === 'right') {
        this.paneAnimationState = 'visibleRight';
      }
      this.setHighZIndexClass = true;
      this.setLowZIndexClass = false;
    } else {
      if (this.slideDirection === 'left') {
        this.paneAnimationState = 'hiddenLeft';
      } else if (this.slideDirection === 'right') {
        this.paneAnimationState = 'hiddenRight';
      }
      setTimeout(() => {
        // This needs to be delayed so we don't see things
        // through the menu when it's still being animated
        this.setLowZIndexClass = true;
        this.setHighZIndexClass = false;
      }, 300);
    }
  }
  public get paneVisible() {
    return this._paneVisible;
  }
  private _paneVisible: boolean = false;
  public paneAnimationState: string = 'hiddenLeft';

  // These classes can be found in common-ux.css
  // I'd like to move them to this component's css file, but the targeting isn't working from this level (yet).
  // Will try to revisit later.
  @HostBinding('class.slide-panel-hide') setHidePanelClass: boolean = true;
  @HostBinding('class.slide-panel-high-z-index') setHighZIndexClass: boolean = false;
  @HostBinding('class.slide-panel-low-z-index') setLowZIndexClass: boolean = false;

  constructor(private panelService: SlidePanelService) {
  }

  ngOnInit() {
    // The panel itself needs to be hidden on page load so users don't see it animating away.
    // For some reason, it starts out displayed before hiding itself.
    setTimeout(() => {
      this.setHidePanelClass = false;
    }, 1000);

    this.panelService.panelVisible.subscribe(
      (res: boolean) => {
        this.paneVisible = res;
      }
    );
  }
}
