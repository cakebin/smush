import { Component, OnInit } from '@angular/core';
import { faGlasses, faFire } from '@fortawesome/free-solid-svg-icons';

@Component({
  selector: 'top-nav-bar',
  templateUrl: './top-nav-bar.component.html',
  styleUrls: ['./top-nav-bar.component.css']
})
export class TopNavBar implements OnInit {
    public faGlasses = faGlasses;
    public faFire = faFire;

    constructor(){
    }

    ngOnInit() {
    }
}
