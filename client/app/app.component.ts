import { Component, OnInit } from '@angular/core';
import { UserManagementService } from './modules/user-management/user-management.service';
import { MatchManagementService } from './modules/match-management/match-management.service';
import { CharacterManagementService } from './modules/character-management/character-management.service';
import { Subscription } from 'rxjs';
import { TagManagementService } from './modules/tag-management/tag-management.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css', './modules/common-ux/common-ux.css']
})
export class AppComponent implements OnInit {
  private matchesLoaded: boolean = false;
  private charactersLoaded: boolean = false;

  constructor(
    private userService: UserManagementService,
    private matchService: MatchManagementService,
    private characterService: CharacterManagementService,
    private tagService: TagManagementService,
  ) {
  }
  ngOnInit() {
    const userSubscription: Subscription = this.userService.cachedUser.subscribe(
      res => {
        if (res && res.isAuthenticated) {
          // We haven't loaded the data yet, so get it once
          if (!this.matchesLoaded && !this.charactersLoaded) {
            this.matchService.loadAllMatches();
            this.characterService.loadAllCharacters();
            this.tagService.loadAllTags();
            this.matchesLoaded = true;
            this.charactersLoaded = true;
          } else {
            // We're done loading data. Kill this subscription
            userSubscription.unsubscribe();
          }
        }
    });
  }

}
