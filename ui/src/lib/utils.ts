import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import { Member } from "..";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function getBaseUrl() {
  if (import.meta.env.DEV) {
    return "http://localhost:8080";
  } else {
    return window.location.origin;
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
  };
}

export const Grades = ["E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8", "E9"];
const AirForceRanks = [
  "AB",
  "Amn",
  "A1C",
  "SrA",
  "SSgt",
  "TSgt",
  "MSgt",
  "SMSgt",
  "CMSgt",
];
const ArmyRanks = [
  "PVT",
  "PV2",
  "PFC",
  "CPL/SFC",
  "SGT",
  "SSG",
  "SFC",
  "MSG",
  "SGM",
];
export function convertGrade(grade: string, service = "F") {
  let ranks: string[] = [];
  switch (service) {
    case "F":
      ranks = AirForceRanks;
      break;
    case "A":
      ranks = ArmyRanks;
  }
  const rank = ranks[Number(grade.at(1)) - 1];
  if (rank === undefined) {
    return "";
  }
  return rank;
}

export function isHigherRank(grade: string, compareTo: string): boolean {
  const gradeNum = Number(grade.at(1));
  const compareToNum = Number(compareTo.at(1));
  return gradeNum > compareToNum;
}
