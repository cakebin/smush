import { Component, OnInit } from '@angular/core';
import { SlidePanelService } from 'client/app/modules/common-ux/components/slide-panel/slide-panel.service';

@Component({
  selector: 'home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  constructor(
    private panelService: SlidePanelService,
  ) {
  }

  ngOnInit() {
  }

  public openLoginPanel() {
    this.panelService.openPanel();
  }
}
