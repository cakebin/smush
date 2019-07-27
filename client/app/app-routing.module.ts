import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MatchesComponent } from './components/matches/matches.component';
import { ProfileViewComponent } from './components/profiles/profile-view.component';
import { ProfileEditComponent } from './components/profiles/profile-edit.component';
import { InsightsComponent } from './components/insights/insights.component';
import { PageNotFoundComponent } from './components/page-not-found/page-not-found.component';
import { HomeComponent } from './components/home/home.component';
import { AuthGuardService as AuthGuard } from './app-auth-guard.service';

const routes: Routes = [
  {
    path: 'home',
    component: HomeComponent,
  },
  {
    path: 'matches',
    component: MatchesComponent,
    canActivate: [AuthGuard],
  },
  {
    path: 'insights',
    component: InsightsComponent,
    canActivate: [AuthGuard],
  },
  {
    path: 'profile',
    canActivate: [AuthGuard],
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
  { path: '',   redirectTo: '/home', pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
