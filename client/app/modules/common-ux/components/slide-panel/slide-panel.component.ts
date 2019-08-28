import { Component, Input, OnInit, HostBinding, HostListener, ElementRef } from '@angular/core';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { SlidePanelService } from './slide-panel.service';

type SlideDirection = 'left' | 'right';
const ANIMATIONLENGTH: number = 300;

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
      transition('* => *', animate(ANIMATIONLENGTH))
    ])
  ]
})
export class SlidePanelComponent implements OnInit {
  @Input() slideDirection: SlideDirection = 'left';

  public set paneVisible(isVisible: boolean) {
    this._paneVisible = isVisible;
    if (isVisible) {
      // Disable click to close for a little bit
      this._canHideByOutsideClick = false;
      if (this.slideDirection === 'left') {
        this.paneAnimationState = 'visibleLeft';
      } else if (this.slideDirection === 'right') {
        this.paneAnimationState = 'visibleRight';
      }
      this._setHighZIndexClass = true;
      this._setLowZIndexClass = false;
      setTimeout(() => {
        // Re-enable click to close when the panel is fully open
        this._canHideByOutsideClick = true;
      }, ANIMATIONLENGTH);
    } else {
      // If you close the menu you can't close it again ...
      this._canHideByOutsideClick = false;
      if (this.slideDirection === 'left') {
        this.paneAnimationState = 'hiddenLeft';
      } else if (this.slideDirection === 'right') {
        this.paneAnimationState = 'hiddenRight';
      }
      setTimeout(() => {
        // This needs to be delayed so we don't see things
        // through the menu when it's still being animated
        this._setLowZIndexClass = true;
        this._setHighZIndexClass = false;
      }, ANIMATIONLENGTH);
    }
  }
  public get paneVisible() {
    return this._paneVisible;
  }
  public paneAnimationState: string = 'hiddenLeft';
  private _paneVisible: boolean = false;
  private _canHideByOutsideClick: boolean = false;

  // These classes can be found in common-ux.css
  // I'd like to move them to this component's css file, but the targeting isn't working from this level (yet).
  // Will try to revisit later.
  @HostBinding('class.slide-panel-hide') private _setHidePanelClass: boolean = true;
  @HostBinding('class.slide-panel-high-z-index') private _setHighZIndexClass: boolean = false;
  @HostBinding('class.slide-panel-low-z-index') private _setLowZIndexClass: boolean = false;

  // Listener for clicks (close on click outside)
  // ***THIS ASSUMES YOU ONLY HAVE ONE PANEL ON THE PAGE.***
  @HostListener('document:click', ['$event.target'])
  public onClick(targetElement): void {
    const clickedInside = this.elementRef.nativeElement.contains(targetElement);
    if (!clickedInside && this._canHideByOutsideClick) {
      this.panelService.closePanel();
    }
  }

  constructor(private elementRef: ElementRef, private panelService: SlidePanelService) {
  }

  ngOnInit() {
    // The panel itself needs to be hidden on page load so users don't see it animating away.
    // For some reason, it starts out displayed before hiding itself.
    setTimeout(() => {
      this._setHidePanelClass = false;
    }, 1000);

    // This is what's doing the actual panel hiding and closing.
    // The setter above should not be explictly called by anything.
    this.panelService.panelVisible.subscribe(
      (res: boolean) => {
        this.paneVisible = res;
      }
    );
  }
}
