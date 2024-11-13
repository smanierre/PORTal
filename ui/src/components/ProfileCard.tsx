import { Member } from "../redux/api";
import { convertGrade } from "../lib/utils";
import { Link } from "react-router-dom";

interface ProfileCardProps {
  className?: string;
  member: Member | null;
  subordinates: Member[];
}

export default function ProfileCard({
  className,
  member,
  subordinates,
}: ProfileCardProps) {
  return (
    <article className={`${className !== undefined ? className : "p-4"}`}>
      <p>
        {convertGrade(member?.grade || "")} {member?.first_name}{" "}
        {member?.last_name}
      </p>
      <p className="py-2">Subordinates:</p>
      <ul className="overflow-scroll">
        {subordinates &&
          subordinates.map((subordinate) => (
            <li key={subordinate.id}>
              <Link to={`/member/${subordinate.id}`}>
                {convertGrade(subordinate.grade)} {subordinate.first_name}{" "}
                {subordinate.last_name}
              </Link>
            </li>
          ))}
      </ul>
    </article>
  );
}
