import { Component, OnInit, Input, Output, EventEmitter } from '@angular/core';
import { Observable } from 'rxjs';
import { debounceTime, distinctUntilChanged, map } from 'rxjs/operators';
import { NgbTypeaheadSelectItemEvent } from '@ng-bootstrap/ng-bootstrap';

const characters = ["Bayonetta", "Bowser", "Bowser Jr.", "Captain Falcon", "Chrom", "Cloud", "Corrin", "Daisy", "Dark Samus", "Diddy Kong", "Donkey Kong", "Dr. Mario", "Duck Hunt", "Falco", "Fox", "Ganondorf", "Greninja", "Ice Climbers", "Ike", "Incineroar", "Inkling", "Jigglypuff", "Joker", "Ken", "King DeDeDe", "King K. Rool", "Kirby", "Link", "Little Mac", "Lucario", "Lucas", "Lucina", "Luigi", "Mario", "Marth", "Mega Man", "Meta Knight", "Mewtwo", "Mii Brawler", "Mii Gunner", "Mii Sword Fighter", "Mr. Game & Watch", "Ness", "Olimar", "Pac-Man", "Palutena", "Peach", "Pichu", "Pikachu", "Pit", "Pokemon Trainer", "Richter", "Ridley", "Rob", "Robin", "Rosalina and Luma", "Roy", "Ryu", "Samus", "Sheik", "Shulk", "Simon", "Snake", "Sonic", "Toon Link", "Villager", "Wario", "Wolf", "Yoshi", "Young Link", "Wii-Fit Trainer", "Zelda", "Zero-Suit Samus"];

@Component({
  selector: 'common-ux-typeahead',
  templateUrl: './typeahead.component.html',
  styleUrls: ['./typeahead.component.css']
})
export class TypeaheadComponent implements OnInit {
  @Input() items:string[] = [];
  @Output() select:EventEmitter<string> = new EventEmitter<string>();

  public userInput: string;

  constructor() { }

  ngOnInit() {
  }

  search = (text$: Observable<string>) =>
    text$.pipe(
      debounceTime(200),
      distinctUntilChanged(),
      map(term => {
        if(term.length < 2){
          return [];
        } 
        else {
          return characters.filter(v => {
            return v.toLowerCase().indexOf(term.toLowerCase()) > -1
          }).slice(0, 10);
        }
      })
    )

  public onBlur(){
    // If the user has cleared the input and blurred out, we need to output a blank value manually
    // because the typeahead does not recognise this as an input "event" per se
    if(this.userInput=="") this.select.emit("");
  }
  public onSelect(eventObject: NgbTypeaheadSelectItemEvent):void {
    this.select.emit(eventObject.item);
  }
  public clear(): void {
    this.userInput = "";
  }

}
