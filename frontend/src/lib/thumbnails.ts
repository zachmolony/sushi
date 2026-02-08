import * as THREE from "three";
import { GLTFLoader } from "three/examples/jsm/loaders/GLTFLoader.js";

const THUMB_SIZE = 256;

let renderer: THREE.WebGLRenderer | null = null;
let scene: THREE.Scene | null = null;
let camera: THREE.PerspectiveCamera | null = null;
let loader: GLTFLoader | null = null;
let webglAvailable: boolean | null = null;

function checkWebGL(): boolean {
  if (webglAvailable !== null) return webglAvailable;
  try {
    const testCanvas = document.createElement("canvas");
    const gl =
      testCanvas.getContext("webgl2") || testCanvas.getContext("webgl");
    webglAvailable = !!gl;
    if (gl) {
      // Clean up test context
      const ext = (gl as WebGLRenderingContext).getExtension(
        "WEBGL_lose_context",
      );
      if (ext) ext.loseContext();
    }
    console.log("[sushi] WebGL available:", webglAvailable);
  } catch {
    webglAvailable = false;
    console.warn("[sushi] WebGL not available");
  }
  return webglAvailable;
}

function ensureRenderer(): boolean {
  if (!checkWebGL()) return false;
  if (renderer) return true;

  try {
    const canvas = document.createElement("canvas");
    canvas.width = THUMB_SIZE;
    canvas.height = THUMB_SIZE;

    renderer = new THREE.WebGLRenderer({
      canvas,
      antialias: true,
      alpha: true,
      preserveDrawingBuffer: true,
    });
    renderer.setSize(THUMB_SIZE, THUMB_SIZE);
    renderer.setClearColor(0x1a1a2e, 1);
    renderer.outputColorSpace = THREE.SRGBColorSpace;
    renderer.toneMapping = THREE.NeutralToneMapping;
    renderer.toneMappingExposure = 1.5;

    scene = new THREE.Scene();
    scene.background = new THREE.Color(0x1a1a2e);

    camera = new THREE.PerspectiveCamera(45, 1, 0.01, 1000);

    // Lighting — bright enough for unlit / vertex-color models
    const ambient = new THREE.AmbientLight(0xffffff, 1.2);
    scene.add(ambient);

    const dirLight = new THREE.DirectionalLight(0xffffff, 1.2);
    dirLight.position.set(5, 10, 7);
    scene.add(dirLight);

    const fillLight = new THREE.DirectionalLight(0xffffff, 0.6);
    fillLight.position.set(-5, 3, -5);
    scene.add(fillLight);

    const backLight = new THREE.DirectionalLight(0xffffff, 0.3);
    backLight.position.set(0, -5, -10);
    scene.add(backLight);

    loader = new GLTFLoader();
    console.log("[sushi] Three.js renderer initialized");
    return true;
  } catch (e) {
    console.error("[sushi] Failed to create WebGL renderer:", e);
    webglAvailable = false;
    return false;
  }
}

export interface ThumbnailResult {
  dataUrl: string;
  polyCount: number;
}

/**
 * Render a GLB/GLTF file to a base64 PNG thumbnail.
 * Expects a URL like "/localfile/?path=..." served by the Go backend.
 * Returns null if WebGL is unavailable.
 */
export async function renderThumbnail(
  url: string,
): Promise<ThumbnailResult | null> {
  if (!ensureRenderer()) {
    console.warn("[sushi] Skipping thumbnail — no WebGL");
    return null;
  }

  // Clean the scene of previous models
  const toRemove: THREE.Object3D[] = [];
  scene!.traverse((child) => {
    if (
      (child instanceof THREE.Mesh || child instanceof THREE.Group) &&
      child.parent === scene
    ) {
      toRemove.push(child);
    }
  });
  toRemove.forEach((obj) => scene!.remove(obj));

  try {
    // Fetch the GLB as an ArrayBuffer
    console.log("[sushi] Fetching:", url);
    const resp = await fetch(url);
    if (!resp.ok) {
      throw new Error(`Fetch failed: HTTP ${resp.status} ${resp.statusText}`);
    }
    const buffer = await resp.arrayBuffer();
    console.log("[sushi] Got buffer:", buffer.byteLength, "bytes");

    if (buffer.byteLength === 0) {
      throw new Error("Empty file");
    }

    // Parse with GLTFLoader
    const gltf = await new Promise<any>((resolve, reject) => {
      loader!.parse(buffer, "", resolve, reject);
    });

    return renderGLTF(gltf);
  } catch (err) {
    console.error("[sushi] Thumbnail render failed for", url, err);
    return null;
  }
}

function renderGLTF(gltf: any): ThumbnailResult {
  const model = gltf.scene;
  scene!.add(model);

  // Count triangles
  let polyCount = 0;
  model.traverse((child: any) => {
    if (child.isMesh && child.geometry) {
      const geom = child.geometry;
      if (geom.index) {
        polyCount += geom.index.count / 3;
      } else if (geom.attributes.position) {
        polyCount += geom.attributes.position.count / 3;
      }
    }
  });
  polyCount = Math.round(polyCount);

  // Auto-frame: compute bounding box and position camera
  const box = new THREE.Box3().setFromObject(model);
  const center = box.getCenter(new THREE.Vector3());
  const size = box.getSize(new THREE.Vector3());
  const maxDim = Math.max(size.x, size.y, size.z);

  // Handle degenerate models (zero size)
  const distance = maxDim > 0 ? maxDim * 1.8 : 5;

  camera!.position.set(
    center.x + distance * 0.7,
    center.y + distance * 0.5,
    center.z + distance * 0.7,
  );
  camera!.lookAt(center);
  camera!.updateProjectionMatrix();

  renderer!.render(scene!, camera!);

  const dataUrl = renderer!.domElement.toDataURL("image/png");

  // Clean up
  scene!.remove(model);
  model.traverse((child: any) => {
    if (child.geometry) child.geometry.dispose();
    if (child.material) {
      if (Array.isArray(child.material)) {
        child.material.forEach((m: any) => m.dispose());
      } else {
        child.material.dispose();
      }
    }
  });

  console.log(
    "[sushi] Rendered thumbnail:",
    dataUrl.length,
    "chars,",
    polyCount,
    "tris",
  );
  return { dataUrl, polyCount };
}
