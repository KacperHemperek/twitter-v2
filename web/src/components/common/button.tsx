import { cva, VariantProps } from "class-variance-authority";
import { cn } from "../../utils/cn";

const button = cva(
  "px-3 py-1 flex items-center justify-center rounded-full transition-all disabled:opacity-50 disabled:cursor-not-allowed lg:py-2 lg:px-4",
  {
    variants: {
      variant: {
        primary: ["bg-blue-500"],
        success: ["bg-green-500"],
        danger: ["bg-red-500"],
      },
    },
    defaultVariants: {
      variant: "primary",
    },
  },
);

type ButtonProps = React.PropsWithChildren<
  {} & React.ButtonHTMLAttributes<HTMLButtonElement> &
    VariantProps<typeof button>
>;

export const TwButton = ({
  children,
  variant,
  className,
  ...props
}: ButtonProps) => {
  return (
    <button className={cn(button({ variant, className }))} {...props}>
      {children}
    </button>
  );
};
