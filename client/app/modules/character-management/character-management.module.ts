import { HttpClientModule } from '@angular/common/http';
import { NgModule, ModuleWithProviders } from '@angular/core';


@NgModule({
  declarations: [],
  imports: [
    HttpClientModule,
  ],
  exports: []
})
export class CharacterManagementModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: CharacterManagementModule,
        providers: [
          {
            provide: 'CharacterApiUrl',
            useValue: '/api/character'
          }
        ]
    };
  }
}
