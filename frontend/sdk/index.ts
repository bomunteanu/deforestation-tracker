import axios from 'axios';
import { isAuthenticated, authToken } from '../src/store'; // Adjust the path if necessary

export interface Area {
  readonly ID?: number;
  AreaName: string;
  TopRightLat: number;
  TopRightLon: number;
  BottomLeftLat: number;
  BottomLeftLon: number;
  DeforestedArea?: number;
  UserID?: number;
}

export interface History {
  readonly ID: number;
  Date?: string;
  ImagePath: string;
  MaskedImagePath: string;
  DeforestedArea: number;
  AreaID: number;
  Area: Area;
}

// Define the structure of the API response
interface ApiResponse<T> {
  data: T;
}

class APIClient {
  private client;

  constructor(baseURL: string) {
    this.client = axios.create({ baseURL });

    // Set Authorization header if token exists
    authToken.subscribe(token => {
      if (token) {
        this.client.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      } else {
        delete this.client.defaults.headers.common['Authorization'];
      }
    });
  }

  // Login method
  async login(username: string, password: string): Promise<void> {
    try {
      const response = await this.client.post<ApiResponse<{ token: string }>>('/login', { username, password });
      const { token } = response.data;

      // Store token in localStorage and update store
      localStorage.setItem('authToken', token);
      authToken.set(token);
      isAuthenticated.set(true);
    } catch (error) {
      console.error('Login failed', error);
      throw error;
    }
  }

  // Logout method
  async logout(): Promise<void> {
    // Clear token from localStorage and update store
    localStorage.removeItem('authToken');
    authToken.set(null);
    isAuthenticated.set(false);
  }

  // Check authentication
  async check(): Promise<any> {
    try {
      const response = await this.client.get<ApiResponse<any>>('/auth/check');
      return response.data;
    } catch (error) {
      console.error('Authentication check failed', error);
      isAuthenticated.set(false);
      throw error;
    }
  }

  // Delete an area
  async deleteArea(areaId: number): Promise<Area> {
    const response = await this.client.delete<ApiResponse<Area>>(`/areas/${areaId}`);
    return response.data.data;
  }

  // Create an area
  async createArea(areaData: any): Promise<Area> {
    const response = await this.client.post<ApiResponse<any>>('/areas', areaData);
    return response.data.data;
  }

  // Get an area by ID
  async getArea(areaId: number): Promise<Area> {
    const response = await this.client.get<ApiResponse<Area>>(`/areas/${areaId}`);
    return response.data.data;
  }

  // Get all areas
  async getAllAreas(): Promise<Area[]> {
    const response = await this.client.get<ApiResponse<Area[]>>('/areas');
    return response.data.data;
  }

  // Get all histories
  async getAllHistories(): Promise<History[]> {
    const response = await this.client.get<ApiResponse<History[]>>('/histories');
    return response.data.data;
  }

  // Get history by ID
  async getHistoryById(historyId: number): Promise<any> {
    const response = await this.client.get<ApiResponse<History>>(`/histories/${historyId}`);
    return response.data;
  }

  async getHistoryByAreaId(areaID: number): Promise<History[]> {
    const response = await this.client.get<ApiResponse<History[]>>(`/histories/area/${areaID}`);
    return response.data.data;
  }

  async getImageByPath(imagePath: string): Promise<any> {
    const response = await this.client.get(`/images/${imagePath}`, {
      responseType: 'blob' // Ensure the response type is blob for image data
    });
    return response.data;
  }
}

export default APIClient;
