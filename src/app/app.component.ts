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
  public showFooterWarnings:boolean = false;
  public warnings: string[] = [];

  constructor(private commonUXService:CommonUXService){
  }
  
  ngOnInit() {
  }

  public createEntry(): void {
    if(!this.validateMatch()){
      // This should never be reached, but in case someone does manage to re-enable the submit button
      // without entering an opponent name... (we'll probably have bigger problems than this, if so)
      this.warnings.forEach(warningMessage => {
        this.commonUXService.showWarningToast(warningMessage);
      });
      return;
    }
    console.log("Saving match:", this.match);
    this.lastSavedMatch = new MatchViewModel();
    this.lastSavedMatch = Object.assign(this.lastSavedMatch, this.match);
    this.resetMatch();

    // Set footer warnings to false so it won't show up until the next mouseenter
    this.showFooterWarnings = false;
  }

  public validateMatch():boolean {
    this.warnings = [];

    if(!this.match.opponentCharacterName){
      this.warnings.push("Opponent character name required.");
    }
    if(!this.match.userCharacterName && this.match.userCharacterGsp){
      this.warnings.push("User GSP must be associated with a user character.");
    }

    if(this.warnings.length) return false;
    else return true;
  }
  private resetMatch(): void {
    this.match = new MatchViewModel(null, null, this.match.userCharacterName, this.match.userCharacterGsp);
    // Need to manually reset the typeahead since it's not a simple input
    this.opponentCharacterNameInput.clear();
  }
}
