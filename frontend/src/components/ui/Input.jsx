import { forwardRef, useState } from 'react';
import { Eye, EyeOff } from 'lucide-react';

export const Input = forwardRef(({ label, error, helperText, icon: Icon, type = 'text', className = '', ...props }, ref) => {
  const [showPassword, setShowPassword] = useState(false);
  const isPassword = type === 'password';
  const inputType = isPassword ? (showPassword ? 'text' : 'password') : type;

  const baseInputClass = `
        w-full px-3 py-2.5 text-sm
        bg-base-100 border border-base-300 rounded-md
        text-base-content placeholder:text-base-content/40
        transition-colors duration-150
        focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20
        disabled:bg-base-200 disabled:cursor-not-allowed disabled:opacity-60
    `;

  const errorClass = error ? 'border-error focus:border-error focus:ring-error/20' : '';

  return (
    <div className="w-full">
      {label && <label className="block text-sm font-medium text-base-content mb-1.5">{label}</label>}

      <div className="relative">
        {Icon && (
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Icon className="h-4 w-4 text-base-content/40" />
          </div>
        )}
        <input
          ref={ref}
          type={inputType}
          className={`${baseInputClass} ${errorClass} ${Icon ? 'pl-10' : ''} ${isPassword ? 'pr-10' : ''} ${className}`}
          {...props}
        />
        {isPassword && (
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute inset-y-0 right-0 pr-3 flex items-center text-base-content/40 hover:text-base-content focus:outline-none transition-colors"
            tabIndex={-1}
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </button>
        )}
      </div>

      {error && <p className="mt-1.5 text-xs text-error">{error}</p>}
      {helperText && !error && <p className="mt-1.5 text-xs text-base-content/50">{helperText}</p>}
    </div>
  );
});

Input.displayName = 'Input';

export const TextArea = forwardRef(({ label, error, helperText, className = '', rows = 4, ...props }, ref) => {
  const baseClass = `
        w-full px-3 py-2.5 text-sm
        bg-base-100 border border-base-300 rounded-md
        text-base-content placeholder:text-base-content/40
        transition-colors duration-150
        focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20
        disabled:bg-base-200 disabled:cursor-not-allowed disabled:opacity-60
        resize-none
    `;

  const errorClass = error ? 'border-error focus:border-error focus:ring-error/20' : '';

  return (
    <div className="w-full">
      {label && <label className="block text-sm font-medium text-base-content mb-1.5">{label}</label>}
      <textarea ref={ref} rows={rows} className={`${baseClass} ${errorClass} ${className}`} {...props} />
      {error && <p className="mt-1.5 text-xs text-error">{error}</p>}
      {helperText && !error && <p className="mt-1.5 text-xs text-base-content/50">{helperText}</p>}
    </div>
  );
});

TextArea.displayName = 'TextArea';

export const Select = forwardRef(({ label, error, helperText, options = [], className = '', placeholder, children, ...props }, ref) => {
  const baseClass = `
        w-full px-3 py-2.5 text-sm
        bg-base-100 border border-base-300 rounded-md
        text-base-content
        transition-colors duration-150
        focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/20
        disabled:bg-base-200 disabled:cursor-not-allowed disabled:opacity-60
        appearance-none cursor-pointer
    `;

  const errorClass = error ? 'border-error focus:border-error focus:ring-error/20' : '';

  return (
    <div className="w-full">
      {label && <label className="block text-sm font-medium text-base-content mb-1.5">{label}</label>}
      <div className="relative">
        <select ref={ref} className={`${baseClass} ${errorClass} ${className} pr-8`} {...props}>
          {placeholder && (
            <option disabled value="">
              {placeholder}
            </option>
          )}
          {options.length > 0
            ? options.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))
            : children}
        </select>
        <div className="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
          <svg className="h-4 w-4 text-base-content/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>
      {error && <p className="mt-1.5 text-xs text-error">{error}</p>}
      {helperText && !error && <p className="mt-1.5 text-xs text-base-content/50">{helperText}</p>}
    </div>
  );
});

Select.displayName = 'Select';
