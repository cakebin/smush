import { Component, Input, Output, EventEmitter } from '@angular/core';

@Component({
  selector: 'common-ux-masked-number-input',
  templateUrl: './masked-number-input.component.html'
})
export class MaskedNumberInputComponent {
  @Input() inputValue: string = '';
  @Output() inputValueChange: EventEmitter<string> = new EventEmitter<string>();

  // This is intended for pageload / account default setting, etc
  public setValue(programmaticValue: number): void {
    if (programmaticValue == null) {
      return;
    }
    this.inputValue = this._getFormattedNumber(programmaticValue.toString());
    this.inputValueChange.emit(this.inputValue);
  }

  public formatNumber(event: KeyboardEvent) {
    // Skip for arrow keys
    if (['Left', 'Right', 'ArrowLeft', 'ArrowRight'].indexOf(event.key) !== -1) {
      return;
    }
    // Format number
    this.inputValue = this._getFormattedNumber(this.inputValue);
    this.inputValueChange.emit(this.inputValue);
  }

  public checkNumber(event: KeyboardEvent) {
    if (event.key === 'e' || event.key === '+' || event.key === '-') {
      return false;
    } else {
      return true;
    }
  }

  private _getFormattedNumber(input: string): string {
    let newValue: string = input;
    newValue = input.replace(/\D/g, '');
    if (newValue !== '') {
      return parseInt(newValue, 10).toLocaleString();
    } else {
      return input;
    }
  }
}
