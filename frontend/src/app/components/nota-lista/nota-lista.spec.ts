import { ComponentFixture, TestBed } from '@angular/core/testing';

import { NotaLista } from './nota-lista';

describe('NotaLista', () => {
  let component: NotaLista;
  let fixture: ComponentFixture<NotaLista>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [NotaLista],
    }).compileComponents();

    fixture = TestBed.createComponent(NotaLista);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
