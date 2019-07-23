import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MatchesComponent } from './components/matches/matches.component';
import { ProfileViewComponent } from './components/profiles/profile-view.component';
import { ProfileEditComponent } from './components/profiles/profile-edit.component';
import { InsightsComponent } from './components/insights/insights.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';

const routes: Routes = [
  {
    path: 'matches',
    component: MatchesComponent,
  },
  {
    path: 'insights',
    component: InsightsComponent,
  },
  {
    path: 'profile',
    children: [
      {
        path: 'view',
        component: ProfileViewComponent,
      },
      {
        path: 'edit',
        component: ProfileEditComponent,
      },
    ]
  },
  { path: '',   redirectTo: '/matches', pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
