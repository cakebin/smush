import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'smush';

  public saveEntry(item: any):void {
    console.log("calling service method to save our new entry!!!!!11 You have selected", item);
  }
}
