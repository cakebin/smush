import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AdminCharacterComponent } from './admin-character.component';

describe('AdminCharacterComponent', () => {
  let component: AdminCharacterComponent;
  let fixture: ComponentFixture<AdminCharacterComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AdminCharacterComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AdminCharacterComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
