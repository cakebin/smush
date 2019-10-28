import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MatchesComponent } from './page-components/matches/matches.component';
import { ProfileComponent } from './page-components/profiles/profile.component';
import { InsightsComponent } from './page-components/insights/insights.component';
import { PageNotFoundComponent } from './page-components/page-not-found/page-not-found.component';
import { HomeComponent } from './page-components/home/home.component';
import { AdminComponent } from './page-components/admin/admin.component';
import { RequestResetComponent } from './page-components/reset-password/request-reset.component';
import { PasswordResetComponent } from './page-components/reset-password/reset.component';
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
        path: 'edit',
        component: ProfileComponent,
      },
    ]
  },
  {
    // No guard for password reset pages
    path: 'reset-password',
    children: [
      {
        path: 'request',
        component: RequestResetComponent,
      },
      {
        path: 'token',
        component: PasswordResetComponent,
      },
    ]
  },
  {
    path: 'admin',
    canActivate: [AuthGuard],
    data: { expectedRole: 'administrator' },
    component: AdminComponent,
  },
  { path: '',   redirectTo: '/home', pathMatch: 'full' },
  { path: '**', component: PageNotFoundComponent }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
