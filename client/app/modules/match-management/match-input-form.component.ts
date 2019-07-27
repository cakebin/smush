import { Component, OnInit, ViewChild, AfterViewInit } from '@angular/core';
import { MatchViewModel, IMatchViewModel, IUserViewModel } from '../../app.view-models';
import { MaskedNumberInputComponent } from '../common-ux/components/masked-number-input/masked-number-input.component';
import { TypeaheadComponent } from '../common-ux/components/typeahead/typeahead.component';
import { CommonUXService } from '../common-ux/common-ux.service';
import { MatchManagementService } from './match-management.service';
import { UserManagementService } from '../user-management/user-management.service';

@Component({
  selector: 'match-input-form',
  templateUrl: './match-input-form.component.html',
})
export class MatchInputFormComponent implements OnInit, AfterViewInit {
  @ViewChild('opponentCharacterNameInput', { static: false }) private opponentCharacterNameInput: TypeaheadComponent;
  @ViewChild('userCharacterNameInput', { static: false }) private userCharacterNameInput: TypeaheadComponent;
  @ViewChild('userCharacterGspInput', { static: false }) private userCharacterGspInput: MaskedNumberInputComponent;

  public match: IMatchViewModel = new MatchViewModel();

  // We have number masking, so these need to be temporarily stored as a string
  public opponentCharacterGspString: string = '';
  public userCharacterGspString: string = '';

  public showFooterWarnings: boolean = false;
  public warnings: string[] = [];
  public isSaving: boolean = false;

  private user: IUserViewModel;

  constructor(
    private commonUXService: CommonUXService,
    private userService: UserManagementService,
    private matchService: MatchManagementService,
    ) {
  }

  ngOnInit() {
    this.userService.cachedUser.subscribe(
      res => {
        if (res) {
          this.user = res;
          this.match.userId = this.user.userId;
          this.match.userCharacterName = res.defaultCharacterName;
          this.userCharacterGspString = res.defaultCharacterGsp.toString();
          this.match.userCharacterGsp = res.defaultCharacterGsp;
          this._formatUserDefaultValues();
        }
    });
  }

  ngAfterViewInit() {
    this._formatUserDefaultValues();
  }

  private _formatUserDefaultValues(): void {
    setTimeout(() => { // Just give in to the dark side
      if (this.userCharacterGspInput && this.match.userCharacterGsp) {
        this.userCharacterGspInput.setValue(this.match.userCharacterGsp);
      }
      if (this.userCharacterNameInput && this.match.userCharacterName) {
        this.userCharacterNameInput.setDefaultValue(this.match.userCharacterName);
      }
    }, 0);
  }

  public createEntry(): void {
    if (!this.validateMatch()) {
      this.warnings.forEach(warningMessage => {
        this.commonUXService.showWarningToast(warningMessage);
      });
      return;
    }

    this.isSaving = true;

    // Transform some data
    if (this.opponentCharacterGspString) {
      this.match.opponentCharacterGsp = parseInt(this.opponentCharacterGspString.replace(/,/g, ''), 10);
    }
    if (this.userCharacterGspString) {
      this.match.userCharacterGsp = parseInt(this.userCharacterGspString.replace(/,/g, ''), 10);
    }
    console.log('Saving match:', this.match);

    this.matchService.createMatch(this.match).subscribe(response => {
      if (response) {
        this.commonUXService.showSuccessToast('Match saved!');
      }
    }, error => {
      this.commonUXService.showDangerToast('Unable to save match.');
    }, () => {
      this.isSaving = false;
    });

    this.resetMatch();

    // Set footer warnings to false so it won't show up until the next mouseenter
    this.showFooterWarnings = false;
  }

  public validateMatch(): boolean {
    this.warnings = [];

    if (!this.match.opponentCharacterName) {
      this.warnings.push('Opponent character name required.');
    }
    if (!this.match.userCharacterName && this.match.userCharacterGsp){
      this.warnings.push('User GSP must be associated with a user character.');
    }
    if (this.warnings.length) {
      return false;
    } else {
      return true;
    }
  }

  private resetMatch(): void {
    this.match = new MatchViewModel(this.user.userId, null, null, null, null, this.match.userCharacterName, this.match.userCharacterGsp);
    // Need to manually mask the user GSP again
    if (this.match.userCharacterGsp) {
      this.userCharacterGspInput.setValue(this.match.userCharacterGsp);
    }
    // Need to manually reset the opponent character typeahead component
    this.opponentCharacterNameInput.clear();
  }
}
