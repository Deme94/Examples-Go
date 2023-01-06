// get (all) with attributes
db.assets.aggregate([
  {
    $lookup: {
      from: "attributes",
      localField: "name",
      foreignField: "metadata.asset_name",
      as: "attributes"
    }
   }
])

// get where, with attributes
db.assets.aggregate([
  {
    $match: {
      "_id": ObjectId("634ac09024b45620c4e2842c")
    }
  },
  {
    $lookup: {
      from: "attributes",
      localField: "name",
      foreignField: "metadata.asset_name",
      as: "attributes",    
    }
   }
])

// get where, with attributes where
db.assets.aggregate([
  {
    $match: {
      "_id": ObjectId("634ac09024b45620c4e2842c")
    }
  },
  {
    $lookup: {
      from: "attributes",
      localField: "name",
      foreignField: "metadata.asset_name",
      pipeline: [
        {
          $match: {
            "_id": ObjectId("635d145ab59ca0fa55ea1c65")
          } 
        }
      ],
      as: "attributes",    
    }
   }
])

// get where, with attributes where (correlated)
db.assets.aggregate([
  {
    $match: {
      "_id": ObjectId("634ac09024b45620c4e2842c")
    }
  },
  {
    $lookup: {
      from: "attributes",
      let: {date: "$date"},     // RENOMBRA LAS VARIABLES DE ASSETS SI VAMOS A COMPARAR CAMPOS DE LAS DOS TABLAS APARTE DE LAS VARIABLES QUE HACEN EL JOIN.
      localField: "name",
      foreignField: "metadata.asset_name",
      pipeline: [
        {
          $match: {
              $expr: { $gte: [ "$$date", "$timestamp" ] }     // OBLIGATORIO USAR $EXPR PARA COMPARAR ENTRE 2 TABLAS
          }
        }
      ],
      as: "attributes"    
    }
   }
])

// get where & project, with attributes where & project
db.assets.aggregate([
  {
    $match: {
      "_id": ObjectId("634ac09024b45620c4e2842c")
    }
  },
  { $project: { "_id": 0, "name": 1 } },
  {
    $lookup: {
      from: "attributes",
      localField: "name",
      foreignField: "metadata.asset_name",
      pipeline: [
        {
          $match: {
            "_id": ObjectId("635d145ab59ca0fa55ea1c65")
          }
        },
        { $project: { "_id": 0, "metadata.name": 1} }
      ],
      as: "attributes",    
    }
   }
])

// get where & project, with attributes where & project + MERGE DE OBJECTS Y POSICIONES DE ARRAY
db.assets.aggregate([
  {
    $match: {
      "_id": ObjectId("634ac09024b45620c4e2842c")
    }
  },
  { $project: { "_id": 0, "name": 1 } },
  {
    $lookup: {
      from: "attributes",
      localField: "name",
      foreignField: "metadata.asset_name",
      pipeline: [
        {
          $match: {
            "_id": ObjectId("635d145ab59ca0fa55ea1c65")
          }
        },
        { $project: { "_id": 0, "metadata.asset_name": 1} }
      ],
      as: "attributes",    
    }
   },
   {
    // Sube el primer attribute devuelto al padre (hace un merge y clona el resultado!)
     $replaceRoot: { newRoot: { $mergeObjects: [ { $arrayElemAt: [ "$attributes", 0 ] }, "$$ROOT" ] } }
   },
   // Mergea el objeto metadata con el padre (el objeto dentro de metadata ahora es parte del padre)
   {
     $replaceRoot: { newRoot: { $mergeObjects: [ "$metadata", "$$ROOT" ] } }
   },
   { $project: { "attributes": 0, "metadata": 0 } } // IMPRESCINDIBLE PARA NO DUPLICAR LA INFORMACION MERGEADA
]) // NOTA MERGES: SI LAS PROPIEDADES QUE SE MERGEAN SE LLAMAN IGUAL QUE LAS DEL PADRE, EL MERGEO QUEDARÁ ENCAPSULADO BAJO EL NOMBRE DEL HIJO, NO SE PODRA MERGEAR EFECTIVAMENTE!