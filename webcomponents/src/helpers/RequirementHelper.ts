export function validateRequirementEditor(f: FormData): boolean {
  if (f.get("type") === "Qualification") {
    if (f.get("qualification_id")) {
      return true;
    }
    return false;
  } else {
    if (!f.get("name")) {
      return false;
    }
    const daysValidFor = f.get("days_valid_for");
    if (!daysValidFor || +daysValidFor === 0) {
      return false;
    }
    if (!f.get("reference")) {
      return false;
    }
    return true;
  }
}
