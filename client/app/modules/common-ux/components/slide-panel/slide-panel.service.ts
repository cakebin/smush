import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { publish, refCount } from 'rxjs/operators';

@Injectable({ providedIn: 'root' })
export class SlidePanelService {
  public panelVisible: BehaviorSubject<boolean> = new BehaviorSubject<boolean>(false);

  public openPanel() {
    console.log('opening panel');
    this.panelVisible.next(true);
    this.panelVisible.pipe(
      publish(),
      refCount()
    );
  }
  public closePanel() {
    console.log('closing panel');
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
