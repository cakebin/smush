import { Directive, HostListener } from '@angular/core';
import { NgControl, NgModel } from '@angular/forms';

@Directive({
  selector: '[number-mask]',
  exportAs:'numberMask'
})
export class NumberMaskDirective {
  constructor(private ngModel: NgModel, private ngControl: NgControl) {
  }

  // This needs to NOT be 'ngModelChange'. Otherwise we end up with stale models!
  @HostListener('keyup', ['$event'])
  onModelChange(event:any) {
    this._onInputChange(event);
  }
  
  public setValue(programmaticValue: number):void {
    if(programmaticValue==null) return;
    
    let maskedValue: string = this._getMaskedValue(programmaticValue.toString());
    if(maskedValue.length){
      this._updateModels(maskedValue);
    }
  }

  private _onInputChange(event:any) {
    let inputValue: string = '';

    if(event.target){
      inputValue = event.target.value;
    }

    let maskedValue: string = this._getMaskedValue(inputValue);
    if(maskedValue.length){
      this._updateModels(maskedValue);
    }
  }
  private _getMaskedValue(inputValue: string): string {
    let newValue:string = inputValue.replace(/\D/g, '');
    if(newValue!=''){
      newValue = parseInt(newValue).toLocaleString();
    }
    return newValue;
  }
  private _updateModels(newValue: string){
    // Updates the Angular model.
    // We can't call this if we use a ngModelChange HostListener cause we get an infinite loop!
    this.ngModel.viewToModelUpdate(newValue);
    //Updates the user display value
    this.ngControl.valueAccessor.writeValue(newValue);
  }
}
