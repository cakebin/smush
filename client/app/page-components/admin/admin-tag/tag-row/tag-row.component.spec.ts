import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TagRowComponent } from './tag-row.component';

describe('TagRowComponent', () => {
  let component: TagRowComponent;
  let fixture: ComponentFixture<TagRowComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TagRowComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TagRowComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
