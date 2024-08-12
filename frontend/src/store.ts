// src/store.ts
import { writable } from 'svelte/store';

export const isAuthenticated = writable(false);
export const authToken = writable<string | null>(localStorage.getItem('authToken'));