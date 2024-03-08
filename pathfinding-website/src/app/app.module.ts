import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { AppComponent } from './app.component';
import { GridComponent } from './grid/grid.component';
import {NgbModule} from "@ng-bootstrap/ng-bootstrap";
import {ControlsComponent} from "./controls/controls.component";

@NgModule({
  declarations: [
    AppComponent,
    GridComponent,
    ControlsComponent
  ],
  imports: [
    BrowserModule,
    NgbModule
  ],
  providers: [GridComponent],
  bootstrap: [AppComponent]
})
export class AppModule { }
