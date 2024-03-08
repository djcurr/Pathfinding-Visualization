import { Component } from '@angular/core';
import {GridComponent} from "./grid/grid.component";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'] // Corrected property name and usage
})
export class AppComponent {
  title = 'Pathfinding Visualization';
  isEraserActive: boolean = false;
  gridSize: number = DefaultGridSize;
  animationSpeed: number = DefaultAnimationSpeed;
  clearGrid: boolean = true;
  startPathfinding: boolean = false;

  constructor(private gridComponent: GridComponent) {  }

  handleEraserToggled(isEraserActive: boolean): void {
    this.isEraserActive = isEraserActive;
  }
  handleAnimationSpeedChanged(newSpeed: number): void {
    this.animationSpeed = newSpeed;
  }

  handleStartPathfindingEvent() {
    this.startPathfinding = !this.startPathfinding;
  }

  handleChangeGridSize(size: number) {
    this.gridSize = size;
  }

  handleClearGrid() {
    this.clearGrid = !this.clearGrid;
  }
}

export const DefaultGridSize: number = 30;
export const DefaultAnimationSpeed: number = 3;
export const AStar: number = 0;
export const Dijkstra: number = 1;
export const DefaultAlgorithm: number = AStar;
export const Algorithms: Map<number, string> = new Map([
  [AStar, "A*"],
  [Dijkstra, "Dijkstra"]
])
