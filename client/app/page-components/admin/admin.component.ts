import { Component, OnInit } from '@angular/core';
import { ICharacterViewModel, IServerResponse } from 'client/app/app.view-models';
import { CharacterManagementService } from 'client/app/modules/character-management/character-management.service';
import { CommonUxService } from 'client/app/modules/common-ux/common-ux.service';

@Component({
  selector: 'admin',
  templateUrl: './admin.component.html',
})
export class AdminComponent implements OnInit {
  public characters: ICharacterViewModel[] = [];
  public editCharacter: ICharacterViewModel = {} as ICharacterViewModel;
  public newCharacter: ICharacterViewModel = {} as ICharacterViewModel;

  constructor(
    private characterService: CharacterManagementService,
    private commonUxService: CommonUxService,
  ) {
  }

  ngOnInit() {
    this.characterService.cachedCharacters.subscribe(
      res => {
        if (res) {
          this.characters = res;
        }
      });
  }

  public onSelectEditCharacter(event: ICharacterViewModel): void {
    if (event) {
      this.editCharacter.characterId = event.characterId;
      // Copy cached character into edit character (we don't want to edit the app version directly)
      const cachedCharacter = this.characters.find(c => c.characterId === event.characterId);
      Object.assign(this.editCharacter, cachedCharacter);
    } else {
      this.editCharacter = {} as ICharacterViewModel;
    }
  }
  public createCharacter(): void {
    if (!this.newCharacter.characterName) {
      // This shouldn't happen unless someone manually re-enables the create button
      this.commonUxService.showWarningToast('Please specify a name for your new character.');
      return;
    }
    this.characterService.createCharacter(this.newCharacter).subscribe(
      (res: IServerResponse) => {
        if (res.success) {
          this.newCharacter = {} as ICharacterViewModel;
          this.commonUxService.showSuccessToast('Character created!');
        } else {
          this.commonUxService.showDangerToast('Unable to create character.');
        }
      },
      error => {
        this.commonUxService.showDangerToast('Unable to create character.');
        console.error(error);
      }
    );
  }
  public updateCharacter(): void {
    if (!this.editCharacter.characterName) {
      // This shouldn't happen unless someone manually re-enables the update button
      this.commonUxService.showWarningToast('Please specify a name for the character.');
      return;
    }
    if (!this._isCharacterNameUnique(this.editCharacter.characterName)) {
      this.commonUxService.showWarningToast('This character already exists.');
      return;
    }

    this.characterService.updateCharacter(this.editCharacter).subscribe(
      (res: IServerResponse) => {
        if (res.success) {
          this.editCharacter = {} as ICharacterViewModel;
          this.commonUxService.showSuccessToast('Character updated!');
        } else {
          this.commonUxService.showDangerToast('Unable to update character.');
        }
      },
      error => {
        this.commonUxService.showDangerToast('Unable to update character.');
        console.error(error);
      }
    );
  }
  private _isCharacterNameUnique(name: string): boolean {
    if (!this.characters) {
      return false;
    }
    return (this.characters.findIndex(c => c.characterName === name) === -1);
  }

}
