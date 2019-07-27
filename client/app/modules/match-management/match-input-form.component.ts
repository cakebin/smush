import { Component, OnInit } from '@angular/core';

import { MatchViewModel, IMatchViewModel, IUserViewModel } from '../../app.view-models';
import { CommonUXService } from '../common-ux/common-ux.service';
import { MatchManagementService } from './match-management.service';
import { UserManagementService } from '../user-management/user-management.service';


@Component({
  selector: 'match-input-form',
  templateUrl: './match-input-form.component.html',
})
export class MatchInputFormComponent implements OnInit {
  public match: IMatchViewModel = new MatchViewModel();
  public characters = this.matchService.characters;

  private _opponentCharacterGspString: string = '';
  private _userCharacterGspString: string = '';
  public set opponentCharacterGspString(value: string) {
    this._opponentCharacterGspString = value;
    this.match.opponentCharacterGsp = parseInt(value.replace(/\D/g, ''), 10);
  }
  public get opponentCharacterGspString(): string {
    return this._opponentCharacterGspString;
  }
  public set userCharacterGspString(value: string) {
    this._userCharacterGspString = value;
    this.match.userCharacterGsp = parseInt(value.replace(/\D/g, ''), 10);
  }
  public get userCharacterGspString(): string {
    return this._userCharacterGspString;
  }

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
        }
    });
  }

  public createEntry(): void {
    if (!this.validateMatch()) {
      this.warnings.forEach(warningMessage => {
        this.commonUXService.showWarningToast(warningMessage);
      });
      return;
    }
    this.isSaving = true;
    console.log('Saving match:', this.match);
    this.matchService.createMatch(this.match).subscribe(response => {
      if (response) {
        this.commonUXService.showSuccessToast('Match saved!');
      }
    }, error => {
      this.commonUXService.showDangerToast('Unable to save match.');
    }, () => {
      this.resetMatch();
      // Set footer warnings to false so it won't show up until the next mouseenter
      this.showFooterWarnings = false;
      this.isSaving = false;
    });
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
    this.match = new MatchViewModel(
      this.user.userId,
      null,
      null,
      null,
      null,
      this.match.userCharacterName,
      this.match.userCharacterGsp
    );
    this.opponentCharacterGspString = '';
  }
}
