import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "../ui/alert-dialog";

interface ChangesPendingAlertProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  accept: () => void;
  cancel: () => void;
}

export default function ChangesPendingAlert({
  open,
  onOpenChange,
  cancel,
  accept,
}: ChangesPendingAlertProps) {
  return (
    <AlertDialog open={open} onOpenChange={onOpenChange}>
      <AlertDialogContent className="bg-white">
        <AlertDialogHeader>
          <AlertDialogTitle>Changes Pending!</AlertDialogTitle>
          <AlertDialogDescription>
            There are changes to this member made. Click Continue Editing to return, or Cancel
            to discard.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel className="bg-red-800 text-white hover:bg-red-300 hover:text-black" onClick={cancel}>Discard</AlertDialogCancel>
          <AlertDialogAction onClick={accept}>Continue Editing</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
