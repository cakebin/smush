import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { publish, refCount } from 'rxjs/operators';

@Injectable({ providedIn: 'root' })
export class SlidePanelService {
  public panelVisible: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  public openPanel() {
    this.panelVisible.next(true);
    this.panelVisible.pipe(
      publish(),
      refCount()
    );
  }
  public closePanel() {
    this.panelVisible.next(false);
    this.panelVisible.pipe(
      publish(),
      refCount()
    );
  }
  public togglePanel() {
    this.panelVisible.next(!this.panelVisible.value);
    this.panelVisible.pipe(
      publish(),
      refCount()
    );
  }
}
