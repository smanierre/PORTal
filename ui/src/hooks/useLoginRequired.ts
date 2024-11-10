import { useNavigate } from "react-router-dom";
import { useGetLoggedInMemberQuery } from "../redux/api";

export default function useLoginRequired() {
  const { data, isLoading, isFetching } = useGetLoggedInMemberQuery();
  const nav = useNavigate();
  if (data === undefined && !isLoading && !isFetching) {
    nav("/login");
  }
}
