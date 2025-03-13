export function validateQualificationEditor(
  f: FormData,
  initial_requirements: Requirement[],
  recurring_requirements: Requirement[],
): boolean {
  const name = f.get("name");
  if (!name) {
    return false;
  }
  const expires = f.get("expires");
  const expiration_interval = f.get("expiration_interval");
  if (
    expires !== null &&
    expiration_interval &&
    +expiration_interval.toString() < 1
  ) {
    return false;
  }
  if (
    initial_requirements.length === 0 &&
    recurring_requirements.length === 0
  ) {
    return false;
  }
  return true;
}
