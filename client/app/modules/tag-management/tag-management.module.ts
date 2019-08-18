import { HttpClientModule } from '@angular/common/http';
import { NgModule, ModuleWithProviders } from '@angular/core';


@NgModule({
  declarations: [],
  imports: [
    HttpClientModule,
  ],
  exports: []
})
export class TagManagementModule {
  static forRoot(): ModuleWithProviders {
    return {
        ngModule: TagManagementModule,
        providers: [
          {
            provide: 'TagApiUrl',
            useValue: '/api/tag'
          }
        ]
    };
  }
}
