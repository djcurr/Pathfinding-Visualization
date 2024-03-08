// src/app/grid/controls.component.ts
import {Component, ElementRef, EventEmitter, OnInit, Output, SimpleChanges, ViewChild} from '@angular/core';
import {WasmService} from "../wasm.service";
import {Algorithms, AStar, DefaultAlgorithm, DefaultAnimationSpeed, DefaultGridSize, Dijkstra} from "../app.component";

@Component({
  selector: 'app-controls',
  templateUrl: './controls.component.html',
  styleUrls: ['./controls.component.css']
})
export class ControlsComponent implements OnInit {
  @Output() eraserToggled = new EventEmitter<boolean>();
  @Output() startPathfindingEvent = new EventEmitter<void>();
  @Output() animationSpeed = new EventEmitter<number>();
  @Output() gridSize = new EventEmitter<number>();
  @Output() clearGridEvent = new EventEmitter<void>();
  @ViewChild('gridSizeSlider') gridSizeSlider: ElementRef | undefined;
  @ViewChild('animationSpeedSlider') animationSpeedSlider: ElementRef | undefined;
  currentSize: number = DefaultGridSize;
  isEraserActive: boolean = false;
  activeAlgorithm: string = "";

  constructor(private wasmService: WasmService) { }

  ngOnInit(): void {
    this.gridSize.emit(DefaultGridSize);
    this.animationSpeed.emit(DefaultAnimationSpeed);
    this.eraserToggled.emit(false)
    this.isEraserActive = false;
    this.currentSize = DefaultGridSize;
    this.setAlgorithm(DefaultAlgorithm);
  }

  resetApplication() {
    if (this.gridSizeSlider && this.animationSpeedSlider) {
      this.gridSizeSlider.nativeElement.value = DefaultGridSize;
      this.animationSpeedSlider.nativeElement.value = DefaultAnimationSpeed;
    }
    this.ngOnInit()
  }

  clearGrid(): void {
    this.clearGridEvent.emit()
  }

  startPathfinding(): void {
    this.startPathfindingEvent.emit();
  }

  changeAnimationSpeed(event: Event): void {
    const target = event.target as HTMLInputElement;
    const speed = Number(target.value);
    this.animationSpeed.emit(speed);
  }

  changeGridSize(event: Event): void {
    const target = event.target as HTMLInputElement;
    const newSize = Number(target.value);
    this.currentSize = newSize;
    this.gridSize.emit(newSize);
  }

  async setAlgorithm(algorithm: number): Promise<void> {
    let algorithmStr = Algorithms.get(algorithm)
    if (algorithmStr) {
      this.activeAlgorithm = algorithmStr;
    }
    await this.wasmService.setActiveAlgorithm(algorithm, this.currentSize, this.currentSize).catch((error) => {
      console.error("Error setting algorithm:", error);
    }).then(() => {
      this.clearGrid()
    });
  }

  toggleEraser(): void {
    this.isEraserActive = !this.isEraserActive;
    this.eraserToggled.emit(this.isEraserActive);
  }

  get algorithmsArray() {
    return Array.from(this.Algorithms.entries());
  }

  protected readonly DefaultGridSize = DefaultGridSize;
  protected readonly DefaultAnimationSpeed = DefaultAnimationSpeed;
  protected readonly AStar = AStar;
  protected readonly Algorithms = Algorithms;
  protected readonly Dijkstra = Dijkstra;
}
