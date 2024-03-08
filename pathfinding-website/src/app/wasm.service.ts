import { Injectable } from '@angular/core';
import {Cell} from "./cell/cell.model";
declare var Go: any;
@Injectable({
  providedIn: 'root',
})

export class WasmService {
  private wasmModule: any;
  private go: any;
  private wasmReady: Promise<void>;
  private resolveWasmReady!: () => void;
  WASM_URL = 'assets/main.wasm';

  constructor() {
    this.go = new Go();
    this.wasmReady = new Promise<void>((resolve) => {
      this.resolveWasmReady = resolve;
    });
    this.initWasm().catch((error) => {
      console.error('WASM could not be loaded: ', error);
    });
  }

  wasmBrowserInstantiate = async (wasmModuleUrl: string, importObject: any) => {
    let response = undefined;

    // Check if the browser supports streaming instantiation
    if (WebAssembly.instantiateStreaming) {
      // Fetch the module, and instantiate it as it is downloading
      response = await WebAssembly.instantiateStreaming(
        fetch(wasmModuleUrl),
        importObject
      );
    } else {
      // Fallback to using fetch to download the entire module
      // And then instantiate the module
      const fetchAndInstantiateTask = async () => {
        const wasmArrayBuffer = await fetch(wasmModuleUrl).then(response =>
          response.arrayBuffer()
        );
        return WebAssembly.instantiate(wasmArrayBuffer, importObject);
      };

      response = await fetchAndInstantiateTask();
    }

    return response;
  };

  runWasmAdd = async () => {
    // Get the importObject from the go instance.
    const importObject = this.go.importObject;
    // Instantiate our wasm module
    this.wasmModule = await this.wasmBrowserInstantiate(this.WASM_URL, importObject);

    // Allow the wasm_exec go instance, bootstrap and execute our wasm module
    this.go.run(this.wasmModule.instance);

  };

  private async initWasm() {
    await this.runWasmAdd();
    console.log('WASM module loaded:', this.wasmModule);
    this.resolveWasmReady();
  }

  public async setStart(x: number, y: number): Promise<boolean> {
    return this.executeWasmFunction('setStart', x, y);
  }

  public async setEnd(x: number, y: number): Promise<boolean> {
    return this.executeWasmFunction('setEnd', x, y);
  }

  public async getStart(): Promise<Point> {
    try {
      const encoded = await this.executeWasmFunction('getStart');

      // For bigint division, ensure you use BigInt literals or constructor for numbers
      const x: bigint = encoded >> BigInt(32); // Shift right to get the upper 32 bits
      const y: bigint = encoded & (BigInt(2**32) - BigInt(1)); // Bitwise AND to get the lower 32 bits

      return new Point(Number(x), Number(y));
    } catch (error) {
      throw new Error("Start node could not be retrieved.");
    }
  }

  public async getEnd(): Promise<Point> {
    try {
      const encoded = await this.executeWasmFunction('getEnd');

      // For bigint division, ensure you use BigInt literals or constructor for numbers
      const x: bigint = encoded >> BigInt(32); // Shift right to get the upper 32 bits
      const y: bigint = encoded & (BigInt(2**32) - BigInt(1)); // Bitwise AND to get the lower 32 bits

      return new Point(Number(x), Number(y));
    } catch (error) {
      throw new Error("End node could not be retrieved.");
    }
  }

  public async setWall(x: number, y: number, isWall: boolean): Promise<boolean> {
    return this.executeWasmFunction('setWall', x, y, isWall);
  }

  public async findPath(): Promise<any> {
    return this.executeWasmFunction('findPath');
  }

  public async clearGrid(): Promise<boolean> {
    return this.executeWasmFunction('clearGrid');
  }

  public async setActiveAlgorithm(name: number, width: number, height: number): Promise<boolean> {
    return this.executeWasmFunction('setActiveAlgorithm', name, width, height);
  }

  public async changeGridSize(width: number, height: number): Promise<boolean> {
    return this.executeWasmFunction('changeGridSize', width, height);
  }

  public async getNumNodes(): Promise<number> {
    return this.executeWasmFunction('getNumNodes');
  }

  public async getNumPathNodes(): Promise<number> {
    return this.executeWasmFunction('getNumPathNodes');
  }

  public async getWidth(): Promise<number> {
    return this.executeWasmFunction('getWidth');
  }

  public async getGrid(): Promise<Cell[][]> {
    const numNodes = await this.getNumNodes().catch((error) => {
      console.error("Error getting numNodes", error);
    })
    const { memory } = this.wasmModule.instance.exports;
    if (!numNodes) {
      throw new Error('No numNodes found');
    }
    const grid = new Int8Array(memory.buffer, 0, numNodes);
    await this.executeWasmFunction('getGrid', grid, numNodes)
    return this.decodeGrid(grid);
  }

  public async getSnapshot(): Promise<Cell[][]> {
    const numNodes = await this.getNumNodes().catch((error) => {
      console.error("Error getting numNodes", error);
    })
    const { memory } = this.wasmModule.instance.exports;
    if (!numNodes) {
      throw new Error('No numNodes found');
    }
    const grid = new Int8Array(memory.buffer, 0, numNodes);
    const status: boolean = await this.executeWasmFunction('getSnapshot', grid, numNodes)
    if (!status) {
      throw new Error('Failed to get snapshots');
    }
    return this.decodeGrid(grid);
  }

  public async getPath(): Promise<Map<String, Point>> {
    const len = await this.getNumPathNodes().catch((error) => {
      console.error("Error getting numNodes", error);
    })
    const { memory } = this.wasmModule.instance.exports;
    if (!len) {
      throw new Error('No numNodes found');
    }
    const path = new Int32Array(memory.buffer, 0, len * 4);
    const status: boolean = await this.executeWasmFunction('getPath', path, len * 4)
    if (!status) {
      throw new Error('Failed to get snapshots');
    }
    return this.decodePath(path);
  }

  private async executeWasmFunction(functionName: string, ...args: any[]): Promise<any> {
    await this.wasmReady;

    if (!this.wasmModule || !(functionName in this.wasmModule.instance.exports)) {
      throw new Error(`${functionName} is not available on the WASM module.`);
    }
    try {
      return await this.wasmModule.instance.exports[functionName](...args);
    } catch (error) {
      console.error(`Error executing ${functionName} in WASM module:`, error);
      throw error;
    }
  }

  public async decodeGrid(encodedGrid: Int8Array): Promise<Cell[][]> {
    const width = await this.getWidth().catch((error) => {
      console.error("Error getting width", error);
    });

    if (!width) {
      throw new Error('No width provided');
    }

    const height = encodedGrid.length / width;
    const grid: Cell[][] = new Array(height);

    for (let y = 0; y < height; y++) {
      grid[y] = new Array(width); // Initialize the row
      for (let x = 0; x < width; x++) {
        const index = y * width + x;
        const encodedNode = encodedGrid[index];

        const visited = (encodedNode & (1 << 0)) !== 0;
        const isWall = (encodedNode & (1 << 1)) !== 0;
        const isStart = (encodedNode & (1 << 2)) !== 0;
        const isEnd = (encodedNode & (1 << 3)) !== 0;
        const isPath = (encodedNode & (1 << 4)) !== 0;

        // Create a new Cell instance with decoded properties
        grid[y][x] = new Cell(x, y, isWall, isStart, isEnd, isPath, visited);
      }
    }
    return grid;
  }

  async decodePath(encodedPath: Int32Array): Promise<Map<String, Point>> {
    const pathMap: Map<String, Point> = new Map();
    const numSegments = encodedPath.length / 4; // 4 values per segment (x1, y1, x2, y2)

    for (let i = 0; i < numSegments; i++) {
      const index = i * 4;
      const startPoint = new Point(encodedPath[index], encodedPath[index + 1]);
      const endPoint = new Point(encodedPath[index + 2], encodedPath[index + 3]);

      // Use the startPoint as the key and endPoint as the value
      // The toString method of Point will be automatically called when it's used as a key
      pathMap.set(startPoint.toString(), endPoint);
    }
    return pathMap;
  }

}

export class Point {
  private readonly x: number;
  private readonly y: number;
  constructor(x: number, y: number) {
    this.x = x;
    this.y = y;
  }

  getX(): number {
    return this.x;
  }

  getY(): number {
    return this.y;
  }

  toString() {
    return `${this.x},${this.y}`; // This will be used as the Map key
  }
}
