import { Directive, HostListener } from '@angular/core';
import { NgControl, NgModel } from '@angular/forms';

@Directive({
  selector: '[numberMask]',
})
export class NumberMaskDirective {
  constructor(private ngModel: NgModel, private ngControl: NgControl) {
  }

  // This needs to NOT be 'ngModelChange'. Otherwise we end up with stale models!
  @HostListener('keyup', ['$event'])
  onModelChange(event:any) {
    this.onInputChange(event);
  }

  onInputChange(event:any) {
    let newVal:string = "";

    if(event.target){
      let input:string = event.target.value;
      newVal = input.replace(/\D/g, '');
      if(newVal!=""){
        newVal = parseInt(newVal).toLocaleString();
      }
    }

    // Updates the Angular model.
    // We can't call this if we use a ngModelChange HostListener cause we get an infinite loop!
    this.ngModel.viewToModelUpdate(newVal);
    //Updates the user display value
    this.ngControl.valueAccessor.writeValue(newVal);
  }
}
