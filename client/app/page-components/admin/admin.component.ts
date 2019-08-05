import { Component, OnInit } from '@angular/core';
import { ICharacterViewModel } from 'client/app/app.view-models';
import { CharacterManagementService } from 'client/app/modules/character-management/character-management.service';

@Component({
  selector: 'admin',
  templateUrl: './admin.component.html',
})
export class AdminComponent implements OnInit {
  public characters: ICharacterViewModel[] = [];
  public editCharacter: ICharacterViewModel = {} as ICharacterViewModel;
  public newCharacter: ICharacterViewModel = {} as ICharacterViewModel;

  constructor(private characterService: CharacterManagementService) {
  }

  ngOnInit() {
    this.characterService.characters.subscribe(
      res => {
        if (res) {
          this.characters = res;
        }
      });
  }

  public onSelectEditCharacter(event: ICharacterViewModel): void {
    if (event) {
      this.editCharacter.characterId = event.characterId;
    }
  }
}
