export interface IInputProps {
  label?: string;
  type?: "text" | "password" | "email" | "number" | "tel";
  placeholder?: string;
  disabled?: boolean;
  required?: boolean;
  loading?: boolean;
  icon?: string;
  error?: string | boolean;
  ui?: any;
}

export interface IBtnProps {
  label?: string;
  type?: "button" | "submit" | "reset";
  block?: boolean;
  loading?: boolean;
  disabled?: boolean;
  to?: string;
  variant?: "solid" | "outline" | "ghost" | "link" | "soft";
  color?:
    | "primary"
    | "error"
    | "secondary"
    | "success"
    | "info"
    | "warning"
    | "neutral";
  icon?: string;
  leadingIcon?: string;
  trailingIcon?: string;
}

export interface ISelectProps {
  label?: string;
  options: T[];
  placeholder?: string;
  searchable?: boolean;
  searchPlaceholder?: string;
  optionAttribute?: string;
  valueAttribute?: string;
  disabled?: boolean;
  error?: string | boolean;
  loading?: boolean;
}

export interface IMultiSelectProps {
  label?: string;
  options: T[];
  placeholder?: string;
  searchable?: boolean;
  optionAttribute?: string;
  valueAttribute?: string;
  disabled?: boolean;
  error?: string | boolean;
  max?: number;
}

export interface IModalProps {
  title?: string;
  description?: string;
  preventClose?: boolean;
  width?: string;
}

interface IBreadcrumbLink {
  label: string;
  to?: string;
  icon?: string;
  click?: () => void;
}

export interface IBreadcrumbProps {
  items: IBreadcrumbLink[];
}
