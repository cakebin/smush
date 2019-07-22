import { Component, OnInit, ViewChild } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { MatchViewModel, IMatchViewModel } from '../../app.view-models';
import { NumberMaskDirective } from '../common-ux/directives/number-mask.directive';
import { TypeaheadComponent } from '../common-ux/components/typeahead/typeahead.component';
import { CommonUXService } from '../common-ux/common-ux.service';
import { MatchManagementService } from './match-management.service';

@Component({
  selector: 'match-input-form',
  templateUrl: './match-input-form.component.html',
})
export class MatchInputFormComponent implements OnInit {
  @ViewChild('opponentCharacterNameInput', { static: false }) private opponentCharacterNameInput: TypeaheadComponent;
  @ViewChild('userCharacterGspInput', { static: false }) private userCharacterGspInput: NumberMaskDirective;

  public match: IMatchViewModel = new MatchViewModel();
  public lastSavedMatch: IMatchViewModel;
  
  public showFooterWarnings:boolean = false;
  public warnings: string[] = [];
  public isSaving: boolean = false;

  constructor(
    private commonUXService:CommonUXService,
    private matchManagementService: MatchManagementService,
    private decimalPipe: DecimalPipe,
    ){
  }
  
  ngOnInit(){

  }

  public createEntry(): void {
    if(!this.validateMatch()){
      this.warnings.forEach(warningMessage => {
        this.commonUXService.showWarningToast(warningMessage);
      });
      return;
    }

    this.isSaving = true;
    
    // Transform some data
    if(this.match.opponentCharacterGsp) {
      this.match.opponentCharacterGsp = parseInt(this.match.opponentCharacterGsp.toString().replace(/,/g, ''));
    }
    if(this.match.userCharacterGsp) {
      this.match.userCharacterGsp = parseInt(this.match.userCharacterGsp.toString().replace(/,/g, ''));
    }
    console.log("Saving match:", this.match);
    this.matchManagementService.createMatch(this.match).subscribe(response => {
      if(response) this.commonUXService.showSuccessToast("Match saved!");
    }, error => {
      this.commonUXService.showDangerToast("Unable to save match.");
    }, () => {
      this.isSaving = false;
    });

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
    this.match = new MatchViewModel(null, null, null, this.match.userCharacterName, this.match.userCharacterGsp);
    // Need to manually mask the user GSP again
    if(this.match.userCharacterGsp) {
      this.userCharacterGspInput.setValue(this.match.userCharacterGsp);
    }
    // Need to manually reset the opponent character typeahead component
    this.opponentCharacterNameInput.clear();
  }
}
