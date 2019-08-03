import { Pipe, PipeTransform } from '@angular/core';


@Pipe({ name: 'stringToNumber' })
export class ToNumberPipe implements PipeTransform {
  transform(value: number | string): number {
    if (value === null) {
      return null;
    }
    if (typeof value === 'string') {
      value = value.replace(/\D/g, '');
      if (value.length) {
        return parseInt(value, 10);
      }
    } else if (typeof value === 'number') {
      return value;
    }
  }
}
