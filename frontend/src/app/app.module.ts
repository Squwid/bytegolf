import { HttpClientModule } from '@angular/common/http';
import { CUSTOM_ELEMENTS_SCHEMA, NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatButtonModule } from '@angular/material/button';
import { MatToolbarModule } from '@angular/material/toolbar';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { RouterModule, Routes } from '@angular/router';
import { NgbModule } from '@ng-bootstrap/ng-bootstrap';
import { AceEditorModule } from 'ng2-ace-editor';
import { AceConfigInterface, ACE_CONFIG } from 'ngx-ace-wrapper';
import { ToastrModule } from 'ngx-toastr';
import { AppComponent } from './app.component';
import { HolesComponent } from './holes/holes.component';
import { NavbarComponent } from './navbar/navbar.component';
import { PlayViewComponent } from './play-view/play-view.component';
import { PreviousComponent } from './play-view/previous/previous.component';

const appRoutes: Routes = [
  { path: '', component: HolesComponent},
  { path: 'hole', component: PlayViewComponent}
];

const DEFAULT_ACE_CONFIG: AceConfigInterface = {

};


@NgModule({
  declarations: [
    AppComponent,
    NavbarComponent,
    HolesComponent,
    PlayViewComponent,
    PreviousComponent
  ],
  imports: [
    RouterModule.forRoot(appRoutes),
    BrowserModule,
    AceEditorModule,
    NgbModule,
    HttpClientModule,
    ToastrModule.forRoot(),
    BrowserAnimationsModule,
    MatToolbarModule,
    MatButtonModule,
    FlexLayoutModule
  ],
  schemas: [
    CUSTOM_ELEMENTS_SCHEMA
  ],
  providers: [
    {
      provide: ACE_CONFIG,
      useValue: DEFAULT_ACE_CONFIG
    }
],
  bootstrap: [AppComponent]
})
export class AppModule { }
