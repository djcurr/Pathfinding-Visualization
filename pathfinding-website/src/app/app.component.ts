import { Component } from '@angular/core';

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
  drawGrid: boolean = false;

  constructor() {  }

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

  handleDrawGrid() {
    this.drawGrid = !this.drawGrid;
  }
}

export const DefaultGridSize: number = 31;
export const DefaultAnimationSpeed: number = 20;
export const AStar: number = 0;
export const Dijkstra: number = 1;
export const DefaultAlgorithm: number = AStar;
export const Algorithms: Map<number, string> = new Map([
  [AStar, "A*"],
  [Dijkstra, "Dijkstra"]
])
