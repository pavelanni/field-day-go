/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./templates/**/*.go.html",
        // add other paths if needed
    ],
    extract: {
        include: ['./templates/**/*.go.html'],
        DEFAULT: content => content.match(/[^<>"'`\\s]*[^<>"'`\\s:]/g) || [],
    },
    safelist: [
        // Add all the classes you use, especially those with colons (md:, hover:, etc.)
        "flex", "flex-col", "md:flex-row", "gap-6", "w-full", "md:w-2/3", "md:w-1/3",
        "items-stretch", "items-center", "justify-center",
        "bg-white", "rounded", "shadow", "p-6", "p-4", "mb-4", "mb-2", "mt-8", "mt-4", "mt-6", "mt-12",
        "text-center", "text-xl", "text-lg", "text-md", "font-bold", "font-semibold", "text-2xl",
        "bg-blue-600", "hover:bg-blue-700", "text-white", "py-2", "px-4", "rounded-md", "rounded-lg",
        "transition", "min-h-screen", "min-h-[60vh]", "max-w-4xl", "max-w-md", "max-w-lg", "mx-auto",
        "overflow-x-auto", "border", "border-b", "even:bg-gray-50", "bg-gray-50", "bg-gray-100",
        "text-gray-500", "text-gray-600", "text-gray-700", "underline", "space-y-4", "space-y-2",
        "block", "focus:outline-none", "focus:ring", "focus:border-blue-400", "focus:ring-2", "focus:ring-blue-500", "focus:ring-offset-2",
        // Add any other classes you see in your templates!
        "md:w-1/3", "lg:col-span-1", "even:bg-gray-50", "lg:flex-row",
        "md:flex-row", "md:grid-cols-3", "hover:bg-blue-700", "lg:w-1/3", "lg:w-2/3", "md:grid-cols-2", "lg:grid-cols-3",
        "focus:outline-none", "focus:ring-2", "focus:ring-blue-500", "focus:ring-offset-2",

    ],
    theme: {
        extend: {},
    },
    plugins: [],
}