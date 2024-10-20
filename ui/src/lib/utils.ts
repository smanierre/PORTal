import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"
import { Member } from ".."

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function getBaseUrl() {
  if (import.meta.env.DEV) {
    return "http://localhost:8080"
  } else {
    return window.location.origin
  }
}

export function getEmptyMember(): Member {
  return {
    first_name: "",
    last_name: "",
    rank: "",
    id: "",
    admin: false,
    username: "",
    supervisor_id: "",
  }
}