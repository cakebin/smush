import { Component, OnInit, ViewChild } from '@angular/core';
import { MatchViewModel, IMatchViewModel } from './app.view-models';
import { TypeaheadComponent } from './modules/common-ux/components/typeahead/typeahead.component';
import { CommonUXService } from './modules/common-ux/common-ux.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  @ViewChild('opponentCharacterNameInput', {static: false}) private opponentCharacterNameInput: TypeaheadComponent;
  public match: IMatchViewModel = new MatchViewModel();
  public lastSavedMatch: IMatchViewModel;
  public warnings: string[] = [];

  constructor(private commonUXService:CommonUXService){
  }
  
  ngOnInit() {
  }

  public createEntry(): void {
    if(!this.match.opponentCharacterName){
      this.commonUXService.showWarningToast("Opponent character name required!");
      return;
    }
    console.log("Saving match:", this.match);
    this.lastSavedMatch = new MatchViewModel();
    this.lastSavedMatch = Object.assign(this.lastSavedMatch, this.match);
    this.resetMatch();
  }

  public test(event:any){
    console.log(event);
  }
  public validateMatch():boolean {
    this.warnings = [];

    if(!(this.match && this.match.opponentCharacterName)){
      this.warnings.push("Opponent character name required.");
      return false;
    }
    else {
      this.warnings = [];
      return true;
    }
  }
  private resetMatch(): void {
    this.match = new MatchViewModel(null, null, this.match.userCharacterName, this.match.userCharacterGsp);
    // Need to manually reset the typeahead since it's not a simple input
    this.opponentCharacterNameInput.clear();
  }
}
