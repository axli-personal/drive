export function getFileViewType(name) {
  let viewType = "binary";
  const dotPos = name.lastIndexOf(".");
  if (dotPos !== -1) {
    switch (name.substring(dotPos + 1)) {
      case "txt":
        viewType = "text"
        break;
      case "md":
        viewType = "markdown"
        break;
    }
  }
  return viewType;
}