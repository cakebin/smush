import { Component, OnInit, ViewChild } from '@angular/core';
import { DecimalPipe } from '@angular/common';
import { IUserViewModel, UserViewModel } from '../../app.view-models';
import { NumberMaskDirective } from '../common-ux/directives/number-mask.directive';
import { TypeaheadComponent } from '../common-ux/components/typeahead/typeahead.component';
import { CommonUXService } from '../common-ux/common-ux.service';
import { UserManagementService } from './user-management.service';

@Component({
  selector: 'user-edit-form',
  templateUrl: './user-edit-form.component.html',
})
export class UserEditFormComponent implements OnInit {
  @ViewChild('defaultCharacterNameInput', { static: false }) private defaultCharacterNameInput: TypeaheadComponent;
  @ViewChild('defaultCharacterGspInput', { static: false }) private defaultCharacterGspInput: NumberMaskDirective;

  public user: IUserViewModel = new UserViewModel();
  
  public showFooterWarnings:boolean = false;
  public warnings: string[] = [];
  public isSaving: boolean = false;

  constructor(
    private commonUXService:CommonUXService,
    private userManagementService: UserManagementService,
    private decimalPipe: DecimalPipe,
    ){
  }
  
  ngOnInit(){

  }
}
